var Client = function(options) {
    var MAX_CONNECT_TIMES = 10;
    var DELAY = 1500;
    var ws;
    this.options = options || {};
    this.createConnect(MAX_CONNECT_TIMES, DELAY);
}

var appendMsg = function(text) {
    var span = document.createElement("SPAN");
    var text = document.createTextNode(text);
    span.appendChild(text);
    document.getElementById("box").appendChild(span);
}

Client.prototype.createConnect = function(max, delay) {
    var self = this;
    if (max === 0) {
        return;
    }
    connect();
    var heartbeatInterval;
    function connect() {
        var urlpath = fmt.Sprintf("ws://localhost:8081/ws?user_id=%d&device_id=%d&token=%s", 22,33, "xxxx")
        self.ws = new WebSocket(urlpath);
        self.ws.binaryType = 'arraybuffer';
        self.ws.onopen = function() {}

        self.ws.onmessage = function(evt) {
            var data = evt.data;
            const buf = new Uint8Array(data);
            const InputMsg = protoRoot.lookup("pb.Input");
            const MessageSync = protoRoot.lookup("pb.SyncResp")
            const responseBody = InputMsg.decode(buf);
            var op = responseBody.type
            switch(op) {
                case 2:
                    var msgBody = MessageSync.decode(responseBody.data);
                    messageReceived(msgBody);
                    var receiveTime = new Date().getTime()
                    appendMsg("batch msg receive: message=" + msgBody+" recieve time:"+receiveTime);
                    break;
                case 3:
                    document.getElementById("status").innerHTML = "<color style='color:green'>ok<color>";
                    appendMsg("receive: auth reply");
                    heartbeat();
                    heartbeatInterval = setInterval(heartbeat, 30 * 1000);
                    break;
                default:
                    console.log("unknow package type=",op)
                    break
            }
        }

        self.ws.onclose = function() {
            if (heartbeatInterval) clearInterval(heartbeatInterval);
            setTimeout(reConnect, delay);

            document.getElementById("status").innerHTML =  "<color style='color:#ff0000'>failed<color>";
        }

        function heartbeat() {
            var buffer = new ArrayBuffer(16); // 创建一个视图，此视图把缓冲内的数据格式化为一个32位（4字节）有符号整数数组
            var int32View = new Int32Array(buffer);
            var payload ={
                type:3,
                requestId:4,
                data:int32View,
            }

            let errMsg = InputMsg.verify(payload);
            if (errMsg != null){
                console.log("buff 解析错误信息：", errMsg);
                return
            }
            // Create a new messageor use .fromObject if conversion is necessary
            const wsData = SignInInput.create(payload);
            // Encode a message to an Uint8Array (browser) or Buffer (node)
            var encoded = SignInInput.encode(wsData).finish();
            self.ws.send(encoded.buffer);
        }

        function messageReceived(body) {
            var notify = self.options.notify;
            if(notify) notify(body);
            console.log("messageReceived body =" + body);
        }
    }

    function reConnect() {
        self.createConnect(--max, delay * 2);
    }
}

((win) =>{
    win['MyClient'] = Client;
    sendMsg()
})(window);
var client = new MyClient({
    notify: function(data) {
        console.log(data);
    }
});

function  sendMsg(){
    var msgBody = document.getElementById("chatMsg")
    var params = {body: msgBody.value,sendTime:new Date().getTime()}
    const SendMessageReq = protoRoot.lookup("pb.SendMessageReq")
    var payload ={
        messageId:1,
        receiverType:2,
        receiverId: 33,
        toUserIds:[2,3],
        sendTime: new Date().getTime(),
        isPersist: true,
    }
    let errMsg = SendMessageReq.verify(payload);
    if (errMsg !=null){
        console.log("buff 解析错误信息：", errMsg);
    }
    const wsData = SendMessageReq.create(payload);
    var encoded = SendMessageReq.encode(wsData).finish();
    console.log(encoded)
    client.ws.send(encoded.buffer)
}