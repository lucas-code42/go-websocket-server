function textMsg() {
    const ws = new WebSocket("ws://localhost:8080/ws-text");
    ws.send("hello from web-socket");
}
textMsg();

function jsonMsg() {
    const ws = new WebSocket("ws://localhost:8080/ws-json?nickname=client-1");
    ws.send(
        JSON.stringify({
            nickname: "client-1",
            msg: "Hello",
            target: "",
            status: "online"
        })
    );
    ws.onmessage = function (event) { console.log(event.data) }
}
jsonMsg();