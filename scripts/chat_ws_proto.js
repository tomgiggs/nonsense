
var protoRoot = require('../proto/proto.js');//由 pbjs -t json-module -w commonjs -o scripts/proto.js pkg/proto/*.proto 生成
const Input = protoRoot.lookup('Input');
const Output = protoRoot.lookup('Output');
const SendMessageReq = protoRoot.lookup('SendMessageReq');

var appendMsg = function (text) {
    var span = document.createElement('SPAN');
    text = document.createTextNode(text);
    span.appendChild(text);
    document.getElementById('box').appendChild(span);
};

function auth() {
}

function heartbeatMsg() {
    return ""
}

function messageReceived( body) {
    console.log('messageReceived: ' + body);
}

const PACKET_TYPE ={
    PT_UNKNOWN:0,
    PT_SIGN_IN : 1, // 设备登录请求
    PT_SYNC:2, // 消息同步触发
    PT_HEARTBEAT:3, // 心跳
    PT_MESSAGE_ACK:4, // 消息投递
    PT_MESSAGE_SEND:5 //正常聊天消息
};
var Client = function (options) {
    var MAX_CONNECT_TIMES = 10;
    var DELAY = 15000;
    this.options = options || {};
    this.createConnect(MAX_CONNECT_TIMES, DELAY);
    var Wsclient;
};

Client.prototype.createConnect = function (max, delay) {
    var self = this;
    if (max === 0) {
        return;
    }
    connect();
    var heartbeatInterval;
    function connect() {
        self.Wsclient = new WebSocket('ws://127.0.0.1:16001/ws');
        self.Wsclient.binaryType = 'arraybuffer';
        self.Wsclient.onopen = function () {
            auth();
            document.getElementById('status').innerHTML ='<color style=\'color:green\'>success<color>';
        };

        self.Wsclient.onmessage = function (evt) {
            var dataView = new Uint8Array(evt.data);
            var decodedMsg = Output.decode(dataView);
            console.log("xxxxxxxxxx",decodedMsg);
            var op = decodedMsg.type;
            switch (op) {
                case PACKET_TYPE.PT_HEARTBEAT:
                    console.log('receive: heartbeat');
                    document.getElementById('status').innerHTML ='<color style=\'color:green\'>ok<color>';
                    ws.send(heartbeatMsg());// send a heartbeat to server
                    heartbeatInterval = setInterval(heartbeat, 30 * 1000);
                    break;
                case PACKET_TYPE.PT_SYNC:
                    var msgBody = SyncResp.decode(decodedMsg.data);
                    console.log("get sync msg:=======",msgBody);
                    messageReceived(msgBody);
                    break;
                default:
                    var msgBody = SyncResp.decode(decodedMsg.data);
                    messageReceived(msgBody);
                    appendMsg('receive: op=' + op + ' message=' + msgBody + ' recieve time:' + new Date().getTime()
                    );
                    break;
            }
        };

        self.Wsclient.onclose = function () {
            if (heartbeatInterval) clearInterval(heartbeatInterval);
            setTimeout(reConnect, delay);
            document.getElementById('status').innerHTML ='<color style=\'color:red\'>failed<color>';
        };
    }
    function reConnect() {
        self.createConnect(--max, delay * 2);
    }
};
Client.prototype.send= function(msg){
    this.Wsclient.send(sendMsg())
}

function sendMsg() {
    var msgReq = SendMessageReq.create({
        messageId:"111",
        receiverType:2,
        receiverId:3,
        toUserId:[3],
        sendTime:new Date().getTime(),
        isPersist:true,
        messageBody:{
            messageType:1,
            messageContent:{
                text:{
                    text:"hello,this is test"
                }
            }
        }
    });
    var inputData = Input.create({
        type:5,
        requestId:"222",
        data:SendMessageReq.encode(msgReq).finish()
    })
    var encoded = Input.encode(inputData).finish();
    console.log("send msg: ",Input.decode(encoded));
    return encoded;
}
window['sendMsg'] = sendMsg;
function main(){
    var client = new Client({
        notify: function (data) {
            console.log(data);
        },
    });
    setTimeout(() => {
        client.send()
    }, 2000);
}
main()

