# protoc-gen-ts-http

## How to use

`protoc --plugin=protoc-gen-ts-http=path/to/plugin --ts-http_out=path/to/out --proto_path=path/to/proto yourProto.proto`

## Codegen demo

`.proto` file:
```proto
service UserService {
  rpc GetUserInfo(GetUserRequest) returns (UserInfo) {
    option(google.api.http) = {
      get: "/web/v1/user/{id}"
    };
  }
}

message GetUserRequest {
  int32 id = 1;
}

message UserInfo {
  int32 id = 1;
  string name = 2;
  optional string employee_no = 3;
  optional int32 department_id = 4;
  optional string department_name = 5;
  optional string job_title = 6; 
}
```

codegen to `.ts` file

```typescript
export class UserApi {
    send: <T = any, R = any>({ method, url, data }: { method: string, url: string, data: T }) => Promise<R>;
    fromRequest: <T = any>(data: T) => any;
    fromResponse: <T = any>(data: T) => any;
    constructor(
        send: <T = any, R = any>({ method, url, data }: { method: string, url: string, data: T }) => Promise<R>,
        fromRequest: <T = any>(data: T) => any,
        fromResponse: <T = any>(data: T) => any,
    ) {
        this.send = send;
        this.fromRequest = fromRequest;
        this.fromResponse = fromResponse;
    }
    
    public async getUserInfo(data: GetUserRequest): Promise<UserInfo> {
        const request_data = this.fromRequest(data)
        return new Promise<UserInfo>((resolve, reject) => {
            
            const id = request_data.id
            this.send({
                method: "GET",
                url: `/web/v1/user/${id}`,
                data: request_data,
            }).then((data) => {
                resolve(this.fromResponse(data) as UserInfo)
            }).catch((error) => {
                reject(error)
            })
        })
    }
}
```

## How to use `.ts` file

```typescript
import {UserApi} from "./useruserApi";

// Constructing the send function
const customSend = async <T, R>({ method, url, data }: { method: string, url: string, data: T }): Promise<R> => {
    const response = await axios({ method, url, data });
    return response.data;
};

const fromRequest = <T = any>(data: T)=> {
    return data
}

const userApi = new UserApi(customSend, fromRequest, fromRequest);

function test(){
    userApi.getUserInfo(new GetUserRequest({id:1})).then((req)=>{
        console.log(req)
    })
}
```