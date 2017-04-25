$(function() {
  startWebSocket();
  autoStopWebSocket();

  $("input[type=radio]").attr('disabled', false);

  $('.answer-option').change(function(event) {
    event.preventDefault();

    var question_number = $(event.target).attr('data-question-number');
    var answer_number = $(event.target).attr('data-answer-number');
    // var data = $().serializeArray();
    // data.push({
    //   selected_answer: 
    // });

    $.ajax({
      url: '/question/' + question_number + '/answer/' + answer_number,
      type: 'POST',
      // data: $.param(data),
    }).done(function(response) {
      console.log(response);
    }).fail(alertAjaxFailure);
  });

  var idxStr = $('#current-question-index').attr('data-question-id');
  if (idxStr != null && parseInt(idxStr) > 0) {
    $('#countdown').text('30 seconds remaining');

    var countDownDate = new Date().getTime() + 30000;
    var x = setInterval(function() {

        var now = new Date().getTime();
        var distance = countDownDate - now;    
        var seconds = Math.floor((distance % (1000 * 60)) / 1000);
        
        if (seconds == 1) {
          $('#countdown').text(seconds + ' second remaining');
        } else {
          $('#countdown').text(seconds + ' seconds remaining');
        }
        
        if (distance < 0) {
          clearInterval(x);
          $('#countdown').text('Time\'s up!');
          $("input[type=radio]").attr('disabled', true);
        }
    }, 1000);    
  } else {
    $('#countdown').text('');
  }
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