package model

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type Part struct {
	UUID      string     `bson:"uuid"`
	Info      PartInfo   `bson:"info"`
	CreatedAt time.Time  `bson:"created_at"`
	UpdatedAt *time.Time `bson:"updated_at, omitempty"`
}

type PartInfo struct {
	Name          string       `bson:"name"`
	Description   string       `bson:"description"`
	Price         float64      `bson:"price"`
	StockQuantity int64        `bson:"stock_quantity"`
	Category      Category     `bson:"category"`
	Dimensions    Dimensions   `bson:"dimensions"`
	Manufacturer  Manufacturer `bson:"manufacturer"`
	Tags          []string     `bson:"tags"`
	Metadata      Metadata     `bson:"metadata"`
}

type Category int

const (
	CategoryEngine Category = iota + 1
	CategoryFuel
	CategoryPorthole
	CategoryWing
)

type Dimensions struct {
	Length float64 `bson:"length"`
	Width  float64 `bson:"width"`
	Height float64 `bson:"height"`
	Weight float64 `bson:"weight"`
}

type Manufacturer struct {
	Name    string `bson:"name"`
	Country string `bson:"country"`
	Website string `bson:"website"`
}

type Metadata map[string]*Value

type Value struct {
	StringValue *string
	Int64Value  *int64
	DoubleValue *float64
	BoolValue   *bool
}

func (m *Metadata) UnmarshalBSON(data []byte) error {
	var raw bson.M
	if err := bson.Unmarshal(data, &raw); err != nil {
		return err
	}

	out := make(Metadata, len(raw))
	for key, v := range raw {
		val, err := bsonPrimitiveToValue(v)
		if err != nil {
			return fmt.Errorf("metadata[%q]: %w", key, err)
		}
		out[key] = val
	}
	*m = out
	return nil
}

func bsonPrimitiveToValue(v any) (*Value, error) {
	switch x := v.(type) {
	case string:
		return &Value{StringValue: &x}, nil
	case bool:
		return &Value{BoolValue: &x}, nil
	case int32:
		i := int64(x)
		return &Value{Int64Value: &i}, nil
	case int64:
		return &Value{Int64Value: &x}, nil
	case float64:
		return &Value{DoubleValue: &x}, nil
	case int:
		i := int64(x)
		return &Value{Int64Value: &i}, nil
	default:
		return nil, fmt.Errorf("unsupported metadata value type %T", v)
	}
}

func (m Metadata) MarshalBSON() ([]byte, error) {
	out := make(bson.M, len(m))
	for key, val := range m {
		if val == nil {
			continue
		}
		switch {
		case val.StringValue != nil:
			out[key] = *val.StringValue
		case val.Int64Value != nil:
			out[key] = *val.Int64Value
		case val.DoubleValue != nil:
			out[key] = *val.DoubleValue
		case val.BoolValue != nil:
			out[key] = *val.BoolValue
		default:
			out[key] = nil
		}
	}
	return bson.Marshal(out)
}
