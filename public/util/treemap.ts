import * as d3 from "d3";
import { Treemap, TreeGraph } from "../store/treemap/index";

export function graphTreegraph(rootTag: string, treegraphData: TreeGraph) {
  var margin = { top: 20, right: 40, bottom: 20, left: 40 };

  const widthRaw = d3.select(rootTag).style("width");
  const width =
    Number(widthRaw.substring(0, widthRaw.length - 2)) -
    margin.left -
    margin.right;

  const heightRaw = d3.select(rootTag).style("height");
  const height =
    Number(heightRaw.substring(0, heightRaw.length - 2)) -
    margin.top -
    margin.bottom;

  // append the svg object to the body of the page
  d3.select(rootTag).selectAll("*").remove();
  var svg = d3.select(rootTag),
    g = svg
      .append("g")
      .attr("transform", "translate(" + margin.left + "," + margin.top + ")");

  var tree = d3.tree().size([height, width - 160]);

  var stratify = d3.stratify().parentId(function (d: any) {
    return d.id.substring(0, d.id.lastIndexOf("."));
  });

  var root = stratify(treegraphData.items).sort(function (a: any, b: any) {
    return a.height - b.height || a.id.localeCompare(b.id);
  });

  var link = g
    .selectAll(".link")
    .data(tree(root).links())
    .enter()
    .append("path")
    .attr("class", "link")
    .attr(
      "d",
      d3
        .linkHorizontal()
        .x(function (d: any) {
          return d.y;
        })
        .y(function (d: any) {
          return d.x;
        }) as any
    );

  var node = g
    .selectAll(".node")
    .data(root.descendants())
    .enter()
    .append("g")
    .attr("class", function (d) {
      return "node" + (d.children ? " node--internal" : " node--leaf");
    })
    .attr("transform", function (d: any) {
      return "translate(" + d.y + "," + d.x + ")";
    });

  node.append("circle").attr("r", 2.5);

  node
    .append("text")
    .attr("dy", 3)
    .attr("x", function (d) {
      return d.children ? -8 : 8;
    })
    .style("text-anchor", function (d) {
      return d.children ? "end" : "start";
    })
    .text(function (d) {
      return d.id.substring(d.id.lastIndexOf(".") + 1);
    });
}

export function graphTreemap(rootTag: string, treemapData: Treemap) {
  var margin = { top: 10, right: 10, bottom: 10, left: 10 };

  const widthRaw = d3.select(rootTag).style("width");
  const width =
    Number(widthRaw.substring(0, widthRaw.length - 2)) -
    margin.left -
    margin.right;

  const heightRaw = d3.select(rootTag).style("height");
  const height =
    Number(heightRaw.substring(0, heightRaw.length - 2)) -
    margin.top -
    margin.bottom;

  // append the svg object to the body of the page
  d3.select(rootTag).selectAll("*").remove();
  var svg = d3
    .select(rootTag)
    .append("g")
    .attr("transform", "translate(" + margin.left + "," + margin.top + ")");

  // Give the data to this cluster layout:
  var root = d3.hierarchy(treemapData).sum(function (d) {
    return d.value;
  });

  // Then d3.treemap computes the position of each element of the hierarchy
  d3.treemap().size([width, height]).padding(2)(root);

  // use this information to add rectangles:
  svg
    .selectAll("rect")
    .data(root.leaves())
    .enter()
    .append("rect")
    .attr("x", function (d: any) {
      return d.x0;
      //return 0;
    })
    .attr("y", function (d: any) {
      return d.y0;
      //return 0;
    })
    .attr("width", function (d: any) {
      return d.x1 - d.x0;
      //return 0;
    })
    .attr("height", function (d: any) {
      return d.y1 - d.y0;
      //return 0;
    })
    .style("stroke", "black")
    .style("fill", "slateblue");

  // and to add the text labels
  svg
    .selectAll("text")
    .data(root.leaves())
    .enter()
    .append("text")
    .attr("x", function (d: any) {
      return d.x0 + 5;
      //return 0;
    }) // +10 to adjust position (more right)
    .attr("y", function (d: any) {
      return d.y0 + 20;
      //return 0;
    }) // +20 to adjust position (lower)
    .text(function (d: any) {
      return d.data.name;
    })
    .attr("font-size", "15px")
    .attr("fill", "white");
}
