package templates

// IndexHTML is ...
const IndexHTML = `
<!DOCTYPE html>
<html>
<head>
    <meta name="viewport" content="width=device-width,initial-scale=1">
    <style type="text/css">
        #main {
            margin: 33vh auto 0;
            transform: translateY(-50%);
            text-align: center;
        }
    </style>
    <title>Extxt</title>
</head>
<body>
    <div id="main">
        <form action="/" method="post" enctype="multipart/form-data">
            <input type="file" name="src_file" accept="image/*" required>
            <button type="submit">text!</button>
        </form>
    </div>
</body>
</html>
`
