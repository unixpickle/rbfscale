# rbfscale

This is an attempt to use [RBF networks](https://en.wikipedia.org/wiki/Radial_basis_function_network) as interpolation for scaling up images. This may seem intractable for large images, but I will attempt to do it efficiently using conjugate gradients.

# Results

Resizing is VERY slow. I resized a 30x36 image to get some comparisons:

<table>
<tr>
<th>Original</th>
<th>Preview.app on OS X</th>
<th>RBF with &sigma;=1</th>
</tr>
<tr>
<td><img src="https://raw.githubusercontent.com/unixpickle/rbfscale/master/samples/input.png"></td>
<td><img src="https://raw.githubusercontent.com/unixpickle/rbfscale/master/samples/preview_osx.png"></td>
<td><img src="https://raw.githubusercontent.com/unixpickle/rbfscale/master/samples/variance_1.png"></td>
</tr>
</table>
