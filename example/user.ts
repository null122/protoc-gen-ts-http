// Code generated by protoc-gen-ts-axios. DO NOT EDIT.
// source: user.proto

import {UserInfo,UserDetailInfo,CreateUserRequest,ChangePasswordRequest,GetUserDetailInfoRequest,GetUserRequest,GetDepartmentUserTableRequest,GetDepartmentUserTableReply,Page} from "./user";


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
    
    
    public async GetUserInfo(data: GetUserRequest): Promise<UserInfo> {
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


    public async GetUserDetailInfo(data: GetUserDetailInfoRequest): Promise<UserDetailInfo> {
        const request_data = this.fromRequest(data)
        return new Promise<UserDetailInfo>((resolve, reject) => {
            
            const id = request_data.id
            this.send({
                method: "GET",
                url: `/web/v1/user/detail/${id}`,
                data: request_data,
            }).then((data) => {
                resolve(this.fromResponse(data) as UserDetailInfo)
            }).catch((error) => {
                reject(error)
            })
        })
    }


    public async GetDepartmentUserTable(data: GetDepartmentUserTableRequest): Promise<GetDepartmentUserTableReply> {
        const request_data = this.fromRequest(data)
        return new Promise<GetDepartmentUserTableReply>((resolve, reject) => {
            
            const department_id = request_data.departmentId
            this.send({
                method: "POST",
                url: `/web/v1/user/table/${department_id}`,
                data: request_data,
            }).then((data) => {
                resolve(this.fromResponse(data) as GetDepartmentUserTableReply)
            }).catch((error) => {
                reject(error)
            })
        })
    }


    public async ChangePassword(data: ChangePasswordRequest): Promise<undefined> {
        const request_data = this.fromRequest(data)
        return new Promise<undefined>((resolve, reject) => {
            
            this.send({
                method: "PUT",
                url: `/web/v1/user/password`,
                data: request_data,
            }).then((data) => {
                resolve(this.fromResponse(data) as undefined)
            }).catch((error) => {
                reject(error)
            })
        })
    }


    public async CreateUser(data: CreateUserRequest): Promise<undefined> {
        const request_data = this.fromRequest(data)
        return new Promise<undefined>((resolve, reject) => {
            
            this.send({
                method: "POST",
                url: `/web/v1/user`,
                data: request_data,
            }).then((data) => {
                resolve(this.fromResponse(data) as undefined)
            }).catch((error) => {
                reject(error)
            })
        })
    }

}