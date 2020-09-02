package main

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

func ConvertDateToGrpcTimestamp(date string) (*timestamp.Timestamp, error) {
	tmpstmp, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return nil, err
	}

	grpcTimestamp, err := ptypes.TimestampProto(tmpstmp)
	if err != nil {
		return nil, err
	}

	return grpcTimestamp, nil
}
