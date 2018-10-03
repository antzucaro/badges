This is badges, a Golang port of the badge-generation code for the free and
open source shooter Xonotic. It was created mainly to simplify the generation
of images containing summary statistics for players of the game to include in
their social media profiles. The original code was written in Python, requiring
its own virtualenv and corresponding symlink to the PyCairo library. This code
attempts to remove that complexity as well as the dependency on Pyramid's
(and XonStat's) bootstrapping functions. By breaking away from these, badge
generation can use more efficient queries and leverage Go's excellent standard
library. Here goes!
