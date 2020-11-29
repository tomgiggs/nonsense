
//下面这段代码莫名其妙就是劳保找不到包：Error: Cannot find module 'google-protobuf'，这个包是全局安装的，代码移动到全局node_module目录下就可以运行了
// var messagePb = require('./message_pb');
// var message = new messagePb.SignInInput() // 创建一个 Input 结构体
// message.setAppId(1)
// message.setDeviceId(2000)
// message.setUserId(1111)
// message.setToken("niceToken")
// var bytes = message.serializeBinary(); //serializeBinary  序列化
// console.log(bytes);

//---------
// const proto = require("protobufjs");
// let model = proto.loadSync('../pkg/proto/message.proto');
// let signInput = model.lookupType('pb.SignInInput'); //msg.Account前面的msg是包名，不是文件名
// let message = signInput.create();
// message.appId=1;
// message.deviceId=2000;
// message.userId=111;
// message.token="";
// let encoded = signInput.encode(message).finish();
// console.log(encoded);

//------------
var protoRoot = require("./proto");
const SignInInput = protoRoot.lookup("pb.SignInInput");

// 编码
function encode() {
    var payload ={
        appId:1,
        deviceId:2,
        token:"ssss",
        userId:3,
    }

    let errMsg = SignInInput.verify(payload);
    if (errMsg !=null){
        console.log("buff 解析错误信息：", errMsg);
    }

    // Create a new message
    const wsData = SignInInput.create(payload); // or use .fromObject if conversion is necessary
    // Encode a message to an Uint8Array (browser) or Buffer (node)
    var encoded = SignInInput.encode(wsData).finish();
    console.log(encoded)
    decode(encoded,(xx)=>{
        console.log("decoded: ",xx)
    })
}

// 解码
function decode(data, cb) {
    const buf = new Uint8Array(data);
    const response = SignInInput.decode(buf);
    // 成功回调
    cb(response);
}
encode()