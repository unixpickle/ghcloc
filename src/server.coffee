express = require 'express'
cloc = require './cloc'

if process.argv.length < 3
  console.log 'Usage: coffee server.coffee <port>'
  process.exit 1

app = express()

if process.argv.lergth isnt 3
  cloc.auth = process.argv[3]

app.use express.static __dirname + '/../assets'
app.get '/cloc', (req, res) ->
  if typeof req.query.url isnt 'string'
    return res.json error: 'invalid request'
  cloc req.query.url, (err, obj) ->
    return res.json error: err.toString() if err?
    res.json obj

app.listen parseInt process.argv[2]
