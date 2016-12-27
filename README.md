# rbfscale

This is an attempt to use [RBF networks](https://en.wikipedia.org/wiki/Radial_basis_function_network) as interpolation for scaling up images. This may seem intractable for large images, but I will attempt to do it efficiently using conjugate gradients.

# Results

Resizing is very slow. On a 2.6GHz Intel Core i7, it takes about 8 seconds to resize a 30x36 image using a variance of 2.

RBF networks don't do much better than conventional interpolation algorithms for image resizing. I resized a 30x36 image to get some comparisons:

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
