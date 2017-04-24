$(function() {
  startWebSocket();
  autoStopWebSocket();
});

function handleQuestionChange() {
  // body...
}


function startWebSocket() {
  console.log("Starting websocket");
  
  socketRocket.start(socketURL).then(function(webSocket) {
    webSocket.onNewMessage = function(message) {
      if (message['next']) {
        var payload = message['next'];
        var question_number = payload['question_number'];

        window.location = "/question/" + question_number;
      } else if (message['show']) {

      } else {
        console.log("Received unhandled websocket message type: " + message);
      }


      var encoderState = message['encoderState'];
      var captionerState = message['captionerState'];

      if (typeof encoderState !== 'undefined') {
        changeEncoderState(encoderState.encoder, encoderState.state);
      } else if (typeof captionerState !== 'undefined') {
        changeCaptionerState(captionerState.captionerId, captionerState.state);
      }      
    };

    webSocket.onerror = function(event) {
      console.log('ERROR: ' + event.data);
    };
  }).catch(function(event) {
    console.log(event);
  });
};

function autoStopWebSocket() {
  $(window).on('beforeunload', function() {
    socketRocket.stop();
  });
};