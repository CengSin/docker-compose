package main

import (
	"context"
	"fmt"
)

type MyCustomServer struct {
}

type StartRequest struct {
	Name string
}

type StartResponse struct {
	Ans string
}

func (s *MyCustomServer) Start(ctx context.Context, req *StartRequest, rsp *StartResponse) error {
	rsp.Ans = fmt.Sprint("Hello ", req.Name)
	return nil
}
