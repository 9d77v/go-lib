syntax = "proto3";

package ${package};
import "github.com/9d77v/go-lib/protos/base.proto";

import "google/protobuf/timestamp.proto";
import "github.com/gogo/protobuf/gogoproto/gogo.proto";

// Enable custom Marshal method.
option (gogoproto.marshaler_all) = true;
// Enable custom Unmarshal method.
option (gogoproto.unmarshaler_all) = true;
// Enable custom Size method (Required by Marshal and Unmarshal).
option (gogoproto.goproto_registration) = true;
// Enable generation of XXX_MessageName methods for grpc-go/status.
option (gogoproto.messagename_all) = true;

service ${Entity}Service {
    rpc Create${Entity}(Create${Entity}Request) returns (Create${Entity}Response) {}
    rpc Update${Entity}(Update${Entity}Request) returns (Update${Entity}Response) {}
    rpc Delete${Entity}(Delete${Entity}Request) returns (Delete${Entity}Response) {}
    rpc Get${Entity}ByID(Get${Entity}ByIDRequest) returns (Get${Entity}ByIDResponse) {}
    rpc List${Entity}(List${Entity}Request) returns (List${Entity}Response) {}
}

message ${Entity}{
    int64 id=1;
    google.protobuf.Timestamp create_time = 2 [
        (gogoproto.stdtime) = true
    ];
    google.protobuf.Timestamp update_time = 3 [
        (gogoproto.stdtime) = true
    ];
}

message Create${Entity}Request {
    ${Entity} ${entity}=1;
}

message Create${Entity}Response{
    protos.Error error=1;
}

message Update${Entity}Request {
    ${Entity} ${entity}=1;
}

message Update${Entity}Response{
    protos.Error error=1;    
}

message Delete${Entity}Request {
    int64 id=1;
}

message Delete${Entity}Response{
    protos.Error error=1;    
}

message Get${Entity}ByIDRequest {
    int64 id=1;
}

message Get${Entity}ByIDResponse{
    protos.Error error=1;
    ${Entity} ${entity}=2;
}

message List${Entity}Request{
    string keyword=1;
    int32 page=2;
    int32 pagesize=3;
}

message List${Entity}Response{
    protos.Error error=1;
    int64 total=2;
    repeated ${Entity} list=3;
}
