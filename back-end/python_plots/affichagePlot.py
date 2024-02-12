import matplotlib.pyplot as plt
import numpy as np


# Ce fichier sert à afficher les plots des points moyens en fonction du niveau de chaque trait de personnalité
# Comment l'utiliser ?

# 1. Lancer localhost:8080/simulate100Championships dans un navigateur
    # Cela permet de simuler d'une traite 100 championnats avec des personnalités aléatoires 
    # pour chacun des pilotes pour chaque championnat
# 2. Lancer localhost:8080/statisticsChampionship dans un navigateur
# 3. Localiser le champ "personnalityAverage" (ctrl+f et prendre le premier résultat, qui contient les données de tous les championnats)
    # Attention, il ne faut pas prendre le champ "personalityAveragePoints" qui contient d'autres types de données
# 4. Copier les données du champ "personnalityAverage" dans la variable personality_average_data ci-dessous (ligne 19)
# 5. Copier le nombre de simulations dans la variable nbSimulations ci-dessous (ligne 20)
    # Le nombre de simulations est indiqué dans les statistiques (champ "nbSimulations")
# 6. Lancer le script python (python3 affichagePlot.py)

# Données du champ personnalityAverage
personality_average_data = ...
nbSimulations = ...

categories = list(personality_average_data.keys())
subcategories = list(personality_average_data[categories[0]].keys())

# Création de 4 plots distincts
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