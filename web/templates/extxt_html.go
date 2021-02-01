package templates

// NOTE: 全面的に変更する。もしかるると要らないかも。
// ExtxtHTML is ...
const ExtxtHTML = `
<!DOCTYPE html>
<html>
<head>
    <meta name="viewport" content="width=device-width,initial-scale=1">
    <style type="text/css">
        #main {
            text-align: center;
        }
        #text p {
            text-align: justify;
        }
        #words p {
            text-align: justify;
        }
    </style>
    <title>Extxt</title>
</head>
<body>
    <div id="main">
        <div id="text">
            <h1>Text</h1>
            <p>
                {{.Text}}
            </p>
        </div>
        <div id="words">
            <h1>Words</h1>
            <p>
            {{range .Words}}
                {{.}} /
            {{end}}
            </p>
        </div>
    </div>
</body>
</html>
`
