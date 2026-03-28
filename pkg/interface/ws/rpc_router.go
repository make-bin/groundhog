package ws

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RPCHandlerFunc is the signature for RPC method handlers.
type RPCHandlerFunc func(ctx context.Context, params json.RawMessage) (interface{}, error)

// rpcRequest is the JSON-RPC style request body.
type rpcRequest struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params,omitempty"`
}

// rpcResponse is the JSON-RPC style success response.
type rpcResponse struct {
	Result interface{} `json:"result"`
}

// rpcErrorResponse is the JSON-RPC style error response.
type rpcErrorResponse struct {
	Error string `json:"error"`
}

// RPCRouter dispatches JSON-RPC style requests to registered handler functions.
type RPCRouter struct {
	CronHandler *CronRPCHandler `inject:""`
	routes      map[string]RPCHandlerFunc
}

// NewRPCRouter creates a new RPCRouter.
func NewRPCRouter() *RPCRouter {
	return &RPCRouter{}
}

// init lazily builds the route table after DI has populated CronHandler.
func (r *RPCRouter) init() {
	if r.routes != nil {
		return
	}
	r.routes = make(map[string]RPCHandlerFunc)
	if r.CronHandler != nil {
		r.routes["cron.list"] = r.CronHandler.HandleCronList
		r.routes["cron.add"] = r.CronHandler.HandleCronAdd
		r.routes["cron.update"] = r.CronHandler.HandleCronUpdate
		r.routes["cron.remove"] = r.CronHandler.HandleCronRemove
		r.routes["cron.run"] = r.CronHandler.HandleCronRun
		r.routes["cron.status"] = r.CronHandler.HandleCronStatus
		r.routes["cron.runs"] = r.CronHandler.HandleCronRuns
	}
}

// RegisterRPCHandler adds a method → handler mapping to the router.
func (r *RPCRouter) RegisterRPCHandler(method string, handler RPCHandlerFunc) {
	r.init()
	r.routes[method] = handler
}

// Handle is the Gin handler for POST /rpc.
// It accepts {"method": "cron.list", "params": {...}} and returns
// {"result": {...}} or {"error": "..."}.
func (r *RPCRouter) Handle(c *gin.Context) {
	r.init()

	var req rpcRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, rpcErrorResponse{Error: "invalid JSON: " + err.Error()})
		return
	}

	if req.Method == "" {
		c.JSON(http.StatusBadRequest, rpcErrorResponse{Error: "method is required"})
		return
	}

	handler, ok := r.routes[req.Method]
	if !ok {
		c.JSON(http.StatusNotFound, rpcErrorResponse{Error: "unknown method: " + req.Method})
		return
	}

	result, err := handler(c.Request.Context(), req.Params)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, rpcErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, rpcResponse{Result: result})
}
