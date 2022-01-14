import * as d3 from "d3";
import { Treemap } from "../store/treemap/index";

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
  var svg = d3
    .select(rootTag)
    .append("svg")
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
