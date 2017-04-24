$(function() {
  startWebSocket();
  autoStopWebSocket();

  $('.answer-option').change(function(event) {
    event.preventDefault();

    var question_number = $(event.target).attr('data-answer-number');
    var data = $().serializeArray();
    data.push({
      selected_answer: $(event.target).attr('data-answer-number')
    });

    $.ajax({
      url: '/question/' + question_number + '/answer',
      type: 'POST',
      dataType: 'json',
      data: $.param(data),
    }).done(function(response) {
      console.log(response);
    }).fail(alertAjaxFailure);
  }); 
});

function startWebSocket() {
  socketRocket.start(socketURL).then(function(webSocket) {
    webSocket.onNewMessage = function(message) {
      var nextMessage = message['next'];
      var showMessage = message['show'];

      if (typeof nextMessage !== 'undefined') {
        var question_number = nextMessage['question_number'];

        window.location = "/question/" + question_number;
      } else if (typeof showMessage !== 'undefined') {

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