# utilization of A* to solve the prefix reversal problem

This is a mini application built in Go to solve the pancake sorting problem.

It's A* package is built to use an interface that will accept any type which has the ability to generate children that are of the same type as itself, and which has a method that allows one to check whether it is a goal node. The A* package also requires a heuristic function be given to it that will estimate the distance to the goal.