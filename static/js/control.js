$(function() {
  $('#submit-next-button').click(function(event) {
    event.preventDefault();

    $.ajax({
      url: '/control/next',
      type: 'POST',
    }).done(function(response) {
    	console.log(response);
    }).fail(alertAjaxFailure);
  });	
});