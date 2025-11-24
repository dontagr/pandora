package transport

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
)

func (h *HTTPManager) PreparationReq(row any) (*bytes.Buffer, error) {
	body, err := h.marshal(row)
	if err != nil {
		return nil, fmt.Errorf("get marshal body: %v", err)
	}

	compressedBody, err := h.compress(body)
	if err != nil {

		return nil, fmt.Errorf("compress: %v", err)
	}

	return compressedBody, nil
}

func (h *HTTPManager) marshal(body any) (*bytes.Buffer, error) {
	modelJSON, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("error marshaling body: %w", err)
	}

	return bytes.NewBuffer(modelJSON), nil
}

func (h *HTTPManager) compress(body *bytes.Buffer) (*bytes.Buffer, error) {
	var compressedBody bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressedBody)
	defer func(gzipWriter *gzip.Writer) {
		err := gzipWriter.Close()
		if err != nil {
			h.log.Errorf("gzipWriter.Close: %v", err)
		}
	}(gzipWriter)

	_, err := gzipWriter.Write(body.Bytes())
	if err != nil {
		return nil, fmt.Errorf("error compressing data: %w", err)
	}
	if err := gzipWriter.Close(); err != nil {
		return nil, fmt.Errorf("error closing Gzip writer: %w", err)
	}

	return &compressedBody, nil
}
