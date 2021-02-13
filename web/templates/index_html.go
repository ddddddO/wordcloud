package templates

// IndexHTML is ...
const IndexHTML = `
<!DOCTYPE html>
<html>
<head>
    <meta name="viewport" content="width=device-width,initial-scale=1">
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">
    <title>WordCloud</title>
</head>
<body>
    <form action="/" method="post" enctype="multipart/form-data">
        <div class="form-group row mx-auto">
            <label for="fileInput" class="col-form-label">Image</label>
            <div class="col-sm-10">
                <input class="form-control-file" id="fileInput" type="file" name="src_file" accept="image/*">
            </div>
        </div>
        <div class="form-group row mx-auto">
            <label for="textareaInput" class="col-form-label">Text</label>
            <div class="col-sm-10">
                <textarea class="form-control" id="textareaInput" name="src_text" rows="5"></textarea>
            </div>
        </div>
        <div class="form-group row mx-auto">
            <button type="submit" class="btn btn-outline-success">Run</button>
        </div>
    </form>
</body>
</html>
`
