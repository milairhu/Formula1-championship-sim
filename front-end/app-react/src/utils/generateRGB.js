export function generateRGB(perso) {
    const r = perso.personality.Aggressivity/5 * 255;
    const g = perso.personality.Concentration/5 * 255;
    const b = perso.personality.Confidence/5 * 255;
    const gr = (1 + perso.personality.Docility) /6 ;
    return `rgba(${r},${g},${b},${gr})`;
}