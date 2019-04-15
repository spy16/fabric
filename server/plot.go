package server

import "text/template"

var plotTemplate = template.Must(template.New("plot").Parse(`
<!DOCTYPE html>
<html>
<meta charset="utf-8">

<head>
  <script src="https://d3js.org/d3.v4.min.js"></script>
  <script src="https://unpkg.com/viz.js@1.8.0/viz.js" type="javascript/worker"></script>
  <script src="https://unpkg.com/d3-graphviz@1.3.1/build/d3-graphviz.min.js"></script>
</head>

<body>
  <div id="graph" style="text-align: center;"></div>

  <script>
    d3.select("#graph")
      .graphviz()
      .renderDot('{{ .graphVizData }}');
  </script>

</body>

</html>
`))
