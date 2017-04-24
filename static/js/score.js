$(function() {
  startWebSocket();
  autoStopWebSocket();
});

function startWebSocket() {
  socketRocket.start(socketURL).then(function(webSocket) {
    webSocket.onNewMessage = function(message) {
      var nextMessage = message['winners'];

      if (typeof nextMessage !== 'undefined') {
        var sessions = nextMessage['sessions'];

        var session = $('#session').attr('data-session-token');

        if ($.inArray(session, sessions) != -1) {
          document.body.style.backgroundColor = "green";
        } else {
          document.body.style.backgroundColor = "red";
        }

        console.log(sessions);
      } else {
        console.log("Received unhandled websocket message type: " + message);        
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