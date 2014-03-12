handleCompletion = (value) ->
  value = value.responseJSON
  $('#result-div').html 'Results: <br /><pre>' + JSON.stringify(value, null, 2) + '</pre>';
  $('#count-button').removeAttr 'disabled'

handleError = ->
  console.log 'got error'
  $('#result-div').html 'ERROR'
  $('#count-button').removeAttr 'disabled'

window.countLines = ->
  $('#count-button').attr 'disabled', 'disabled'
  comp = encodeURIComponent $('#repository').val()
  $.ajax "/cloc?url=#{comp}",
    complete: handleCompletion
    error: handleError
    dataType: 'json'
