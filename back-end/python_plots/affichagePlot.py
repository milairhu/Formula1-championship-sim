import matplotlib.pyplot as plt
import numpy as np

# This file is used to display the plots of average points based on the level of each personality trait
# How to use it?

# 1. Run localhost:8080/simulate100Championships in a browser
    # This allows you to simulate 100 championships in one go with random personalities 
    # for each of the drivers for each championship
# 2. Run localhost:8080/statisticsChampionship in a browser
# 3. Locate the "personnalityAverage" field (ctrl+f and take the first result, which contains the data from all the championships)
    # Be careful, do not take the "personalityAveragePoints" field which contains other types of data
# 4. Copy the data from the "personnalityAverage" field into the personality_average_data variable below (line 19)
# 5. Copy the number of simulations into the nbSimulations variable below (line 20)
    # The number of simulations is indicated in the statistics (field "nbSimulations")
# 6. Run the python script (python3 affichagePlot.py)

# Data from the personnalityAverage field
personality_average_data = ...
nbSimulations = ...

categories = list(personality_average_data.keys())
subcategories = list(personality_average_data[categories[0]].keys())

# Creation of 4 distinct plots
fig, axes = plt.subplots(nrows=2, ncols=2, figsize=(10, 8))
fig.suptitle('Personality Traits Average Values')

for i, ax in enumerate(axes.flat):
    category = categories[i]
    values = [personality_average_data[category][subcat]/nbSimulations for subcat in subcategories]
    ax.plot(subcategories, values, alpha=0.7)
    ax.set_title(category)
    ax.set_xlabel('Subcategories')
    ax.set_ylabel('Average Values')

plt.tight_layout(rect=[0, 0.03, 1, 0.95])
plt.show()