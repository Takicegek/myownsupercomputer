$(document).ready(function(){
  var ws;
  var open = false;
  var work = {};

  connectWS = function () {
    console.log("Attempting to connect.");

    ws = new WebSocket("ws://{{ .ServerHost }}:{{ .ServerPort }}/api/ws");

    ws.onclose = onCloseWS;
    ws.onerror = onErrorWS;
    ws.onopen = onOpenWS;
    ws.onmessage = onMessageWS;
  }

  onCloseWS = function () {
    ws = NULL;
    open = false;

    console.log("WebSocket closed.  Attempting to reconnect...");

    // reconnect after 1 second
    setTimeout(function() {connectWS();}, 1000);
  }

  onOpenWS = function () {
    console.log("WebSocket Open.");

    open = true;

    setTimeout(askForWork, 0);
  }

  onErrorWS = function (err) {
    console.log("WebSocket Error: ", err);

    // unsure of how to proceed
  }

  onMessageWS = function (event) {
    console.log("Message Received: ", event.type, event.data);
    switch(event.type) {
      case "message":
        work = JSON.parse(event.data);
        setTimeout(startWork,0);
        break;
      default:
        console.log("Unknown Message Type: ", event.type);
        console.log("   Msg Data: ", event.data);
        break;
    }
  }

  startWork = function () {
      console.log("Starting Work...");
      if (work == {}) {
        return;
      }

      if (work.js != undefined) {
        eval(work.js);
      }
    }

  stopWork = function () {
      console.log("Stopping Work...");
      work = {};
    }

  askForWork = function () {
      console.log("Asking For Work...");
      ws.send(JSON.stringify({type: "need", text: "", id: "{{ .NodeID }}", date: Date.now()}));
    }

   connectWS();

   console.log("Starting Up!");
});
