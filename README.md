Precomputed shade FOV algo implementation
===


This is a implementation of an elegant FOV algorithm for rogouelike games in Golang.
Take a look for (orifinal post)[https://www.reddit.com/r/roguelikedev/comments/5n1tx3/fov_algorithm_sharencompare/] for details.

The problem of this particular algo is walls lighting. My solution is far from ideal, adds much complexity to a
single-pass (in ideal conditions) algo.


Basic description
---

####Method

#####Beforehand

- List the cells in your largest-possible FOV, storing X and Y values relative to the center.
- Store the distance from the center for each cell, and sort the list by this in ascending order.
- Store the range of angles occludedAngles by each cell in this list, in clockwise order as absolute integers only.
- Create a 360-char string of 0s called EmptyShade, and a 360-char string of 1s called FullShade

#####Runtime

- Store two strings – CurrentShade and NextShade
- Set CurrentShade to EmptyShade to start.
- While CurrentShade =/= FullShade: step through the Cell List:

	- If the distance to the current cell is not equal to the previous distance checked then replace the contents
      of the CurrentShade variable with the contents of the NextShade variable.

	- If the tested cell is opaque – for each angle in the range occludedAngles by the cell, place a 1 at the position
      determined by angle%360 in the NextShade string.

    - For each angle in the range occludedAngles by the cell, add 1 to the lit value for that cell for each 0
      encountered at the position determined by angle%360 in the CurrentShade string.

Usage
--
See test for usage example. The types package is a minimal sample, you can use your own,
but edit the algo accordingly.