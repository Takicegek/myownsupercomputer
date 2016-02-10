  var StartNonce = {{ .StartNonce }};
  var EndNonce   = {{ .EndNonce }};
  var Payload    = "{{ .Payload }}";

  for (var i = StartNonce; i != EndNonce; i++) {
    var compute = Sha256.hash(Payload+i);
  }

  ws.send(JSON.stringify({type:"payload", payload: "50000 Hashes Done!"}));

  askForWork();
