// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {main} from '../models';
import {grpc} from '../models';

export function StopRequest(arg1:number):Promise<Error>;

export function OpenImportPath():Promise<string|Error>;

export function OpenProtoFile():Promise<main.OpenProtoFileResult|Error>;

export function RefreshProtoDescriptors(arg1:Array<string>,arg2:Array<string>):Promise<Array<grpc.ProtoTreeNode>|Error>;

export function SelectMethod(arg1:string):Promise<string|Error>;

export function SendRequest(arg1:number,arg2:string,arg3:string,arg4:string):Promise<string|Error>;
