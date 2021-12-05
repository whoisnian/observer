package server

var generatedHtml []byte = []byte(`<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>observer</title>
  </head>
  <body>
    <img src="/stream">
    <script>
    window.isCombo = false;
    document.addEventListener('keydown', (e) => {
      console.log(e)
      e.preventDefault()
      e.stopPropagation()

      url = window.isCombo ? '/api/event?combo=true' : '/api/event'
      window.isCombo = false;
      fetch(url, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          Ctrl: e.ctrlKey,
          Alt: e.altKey,
          Shift: e.shiftKey,
          Code: e.code
        })
      })
      .then(res => res.json())
      .then(data => window.isCombo = data.Combo)
    })
    </script>
  </body>
</html>`)
