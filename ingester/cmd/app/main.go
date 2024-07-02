package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	mygrpc "github.com/nnurry/gopds/ingester/internal/api/grpc"
	"github.com/nnurry/gopds/ingester/internal/database/postgres"
	pb "github.com/nnurry/gopds/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/itchyny/gojq"
)

type CardinalConfig struct {
	Type string `json:"type"`
}

type FilterConfig struct {
	Type           string  `json:"type"`
	MaxCardinality uint    `json:"max_cardinality"`
	ErrorRate      float64 `json:"error_rate"`
}

type IngestConfig struct {
	CardinalConfig CardinalConfig `json:"cardinal_config"`
	FilterConfig   FilterConfig   `json:"filter_config"`
	KeyJqPath      string         `json:"key_jq_path"`
	ValueJqPath    string         `json:"value_jq_path"`
}

func loadRaw(body io.Reader) []byte {
	js, err := io.ReadAll(body)
	if err != nil {
		panic(err)
	}
	return js
}

func loadJson(body io.Reader, dest interface{}) {
	err := json.Unmarshal(loadRaw(body), &dest)
	if err != nil {
		panic(err)
	}
}

func loadJsonFromRaw(raw []byte, dest interface{}) {
	err := json.Unmarshal(raw, &dest)
	if err != nil {
		panic(err)
	}
}

var GrpcProtoClient pb.BatchIngestClient
var IngestConfigs = make(map[string]*IngestConfig)

func loadConfig(configId string) (*IngestConfig, error) {
	var jsonCardinal []byte
	var jsonFilter []byte
	cfg := &IngestConfig{}
	intConfigId, err := strconv.Atoi(configId)
	if err != nil {
		panic(err)
	}
	err = postgres.Client.QueryRow(`
	SELECT key_path, value_path, cardinal_config, filter_config
	FROM json_config
	WHERE id = $1
	`,
		intConfigId,
	).Scan(
		&cfg.KeyJqPath, &cfg.ValueJqPath,
		&jsonCardinal, &jsonFilter,
	)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonCardinal, &cfg.CardinalConfig)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonFilter, &cfg.FilterConfig)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func getConfig(configId string) *IngestConfig {
	var err error
	var cfg *IngestConfig
	cfg, ok := IngestConfigs[configId]
	if !ok {
		cfg, err = loadConfig(configId)
		if err != nil {
			panic(err)
		}
	}
	return cfg
}

func jqFind(jqPath string, input interface{}) string {
	query, err := gojq.Parse(jqPath)
	if err != nil {
		panic(err)
	}
	var value string
	iter := query.Run(input)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			if err, ok := err.(*gojq.HaltError); ok && err.Value() == nil {
				break
			}
			log.Fatalln(err)
		}
		if value == "" {
			value = v.(string)
			break
		}
	}
	return value
}

func CreateConfig(w http.ResponseWriter, r *http.Request) {
	body := &IngestConfig{}

	loadJson(r.Body, &body)

	cc, err := json.Marshal(body.CardinalConfig)
	if err != nil {
		panic(err)
	}
	fc, err := json.Marshal(body.FilterConfig)
	if err != nil {
		panic(err)
	}

	var configId int

	tx, _ := postgres.Client.Begin()

	err = postgres.Client.QueryRow(`
	INSERT INTO json_config (key_path, value_path, cardinal_config, filter_config) 
	VALUES ($1, $2, $3, $4)
	RETURNING id
	;`,
		&body.KeyJqPath, &body.ValueJqPath,
		cc, fc,
	).Scan(&configId)

	if err != nil {
		log.Println("Can't insert JSON config:", err)
		tx.Rollback()
		w.Write([]byte("Failed: " + err.Error()))
		return
	}

	strConfigId := strconv.Itoa(configId)

	IngestConfigs[strConfigId] = body

	tx.Commit()
	w.Write([]byte("Created config of ID = " + strConfigId))
}

func Ingest(w http.ResponseWriter, r *http.Request) {
	var rawData interface{}
	configId := r.PathValue("configId")
	if configId == "" {
		panic(errors.New("invalid config ID: " + configId))
	}
	// NOTE: no schema validation code for input payload yet
	byteData := loadRaw(r.Body)

	_, err := postgres.Client.Exec(`INSERT INTO raw_data (raw_json) VALUES ($1)`, byteData)
	loadJsonFromRaw(byteData, &rawData)
	if err != nil {
		panic(err)
	}
	cfg := getConfig(configId)

	var probKey string
	var probValue string

	probKey = jqFind(cfg.KeyJqPath, rawData)
	probValue = jqFind(cfg.ValueJqPath, rawData)

	if probKey == "" || probValue == "" {
		panic(errors.New("missing key or value: " + probKey + ", " + probValue))
	}

	cardinalType, ok := pb.CardinalType_value[cfg.CardinalConfig.Type]
	if !ok {
		panic("Can't find " + cfg.CardinalConfig.Type + " in enum set")
	}

	filterType, ok := pb.FilterType_value[cfg.FilterConfig.Type]
	if !ok {
		panic("Can't find " + cfg.FilterConfig.Type + " in enum set")
	}

	ingestRequest := &pb.IngestRequest{
		Meta: &pb.MetaField{
			UtcNow: timestamppb.New(time.Now().UTC()),
			Key:    probKey,
			Value:  probValue,
		},
		Cardinal: &pb.CardinalField{
			Type: pb.CardinalType(cardinalType),
		},
		Filter: &pb.FilterField{
			Type:           pb.FilterType(filterType),
			MaxCardinality: uint32(cfg.FilterConfig.MaxCardinality),
			ErrorRate:      float32(cfg.FilterConfig.ErrorRate),
		},
	}
	mygrpc.IngestChannel <- ingestRequest
	w.Write([]byte("Success"))
}

func main() {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(":50051", opts...)
	if err != nil {
		log.Fatal("Can't create gRPC client", err)
	}

	defer conn.Close()

	postgres.Bootstrap()

	GrpcProtoClient = pb.NewBatchIngestClient(conn)

	go func() {
		mygrpc.BatchIngest(GrpcProtoClient)
	}()

	mux := http.NewServeMux()

	mux.HandleFunc("/ingest/config", CreateConfig)
	mux.HandleFunc("/ingest/add/{configId}", Ingest)

	http.ListenAndServe(":6000", mux)
}
