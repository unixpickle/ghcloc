https = require 'https'
url = require 'url'
sloc = require 'sloc'
path = require 'path'

doRawCount = (contents) ->
  lines = contents.split '\n'
  totalCount = lines.length
  sourceLines = 0
  for x in lines
    for c in x
      if not c in [' ', '\t', '\r']
        sourceLines++
        break
  result =
    loc: totalCount
    sloc: sourceLines
    cloc: 0
    scloc: 0
    mcloc: 0
    nloc: 0
  return result

sourceCountLines = (filename, contents) ->
  languageNames =
    c: 'C'
    cc: 'C++'
    cpp: 'C++'
    'c++': 'C++'
    js: 'JavaScript'
    java: 'Java'
    s: 'Assembly'
    asm: 'Assembly'
    go: 'Go'
    cs: 'C#'
    coffee: 'CoffeeScript'
    coffeescript: 'CoffeeScript'
    m: 'Objective-C'
    h: 'C Header'
    md: 'Markdown'
    makefile: 'Makefile'
    mk: 'Makefile'
  
  extname = path.extname(filename).toLowerCase()
  extname = filename.toLowerCase() if extname is ''
  extname = '.c' if extname is '.h'
  extname = extname[1..] if extname[0] is '.'
  
  try
    result = sloc contents, extname
  catch e
    result = doRawCount contents
  
  language = languageNames[extname] ? 'Other'
  ret = new Object()
  ret[language] = result
  return ret

class Repository
  constructor: (@user, @repo) ->
  
  getContentsURL: (aPath) ->
    "https://api.github.com/repos/#{@user}/#{@repo}/contents/#{aPath}"
  
  cloc: (aPath, cb) ->
    req = url.parse @getContentsURL aPath
    req.headers = 'user-agent': 'cloc0.0.1'
    req.auth = module.exports.auth
    x = https.get req, (res) =>
      buff = new Buffer 0
      res.on 'data', (d) -> buff = Buffer.concat [buff, d]
      res.on 'end', =>
        try
          object = JSON.parse buff
        catch e
          return cb e
        if object instanceof Array
          @_clocDirectory object, cb
        else @_clocFile object, cb
    x.on 'error', (error) -> cb error
  
  _clocDirectory: (listing, cb, result = {}) ->
    return cb null, result if listing.length is 0
    [obj, listing] = [listing[0], listing[1..]]
    @cloc obj.path, (err, dict) =>
      return cb err if err?
      for key, obj of dict
        if result[key]?
          for val in ['loc', 'sloc', 'cloc', 'scloc', 'mcloc', 'nloc']
            result[key][val] += obj[val]
        else result[key] = obj
      @_clocDirectory listing, cb, result
  
  _clocFile: (file, cb) ->
    content = new Buffer(file.content, 'base64').toString()
    result = sourceCountLines file.name, content
    cb null, result

module.exports = (url, cb) ->
  parts = /github.com\/(.*?)\/(.*?)($|\/)/.exec url
  return cb new Error 'invalid URL' if not parts
  repo = new Repository parts[1], parts[2]
  repo.cloc '', cb
