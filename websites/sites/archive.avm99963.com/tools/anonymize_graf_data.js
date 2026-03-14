const fs = require("fs");

// IDs that define regions/boundaries and should NOT be touched.
const BLOCKLIST = [205, 206, 207, 208, 559, 560, 562, 563, 10235, 10236, 10237];

const ADJECTIVES = [
  "Quantum",
  "Sleepy",
  "Brave",
  "Cosmic",
  "Silent",
  "Neon",
  "Cyber",
  "Magic",
  "Lost",
  "Crimson",
  "Azure",
  "Golden",
  "Wandering",
  "Hidden",
  "Rapid",
  "Lunar",
  "Solar",
  "Electric",
  "Galactic",
  "Stellar",
  "Crystal",
  "Velvet",
  "Iron",
  "Silver",
  "Jade",
  "Amber",
  "Shadow",
  "Frost",
  "Flame",
  "Storm",
  "Mystic",
  "Astral",
  "Radiant",
  "Clever",
  "Swift",
  "Fierce",
  "Noble",
  "Ancient",
  "Digital",
  "Infinite",
  "Secret",
  "Bold",
  "Calm",
  "Wild",
  "Crisp",
  "Smooth",
  "Hollow",
  "Vibrant",
  "Ethereal",
  "Mighty",
];

const NOUNS = [
  "Badger",
  "Penguin",
  "Fox",
  "Dragon",
  "Phoenix",
  "Wolf",
  "Owl",
  "Bear",
  "Raven",
  "Tiger",
  "Panther",
  "Falcon",
  "Mantis",
  "Cobra",
  "Griffin",
  "Kraken",
  "Sphinx",
  "Yeti",
  "Bison",
  "Leopard",
  "Eagle",
  "Hawk",
  "Shark",
  "Whale",
  "Dolphin",
  "Octopus",
  "Spider",
  "Scorpion",
  "Rhino",
  "Moose",
  "Rabbit",
  "Tortoise",
  "Turtle",
  "Frog",
  "Viper",
  "Pegasus",
  "Unicorn",
  "Gargoyle",
  "Golem",
  "Comet",
  "Asteroid",
  "Meteor",
  "Nebula",
  "Pulsar",
  "Quasar",
  "Galaxy",
  "Planet",
  "Moon",
  "Star",
  "Phantom",
];

const STDIN_FILE_DESCRIPTOR = 0;

const usedFakeNames = new Set();

function run() {
  try {
    const rawData = fs.readFileSync(STDIN_FILE_DESCRIPTOR, "utf8");

    if (!rawData) {
      console.error("Error: No data received on STDIN.");
      process.exit(1);
    }

    const graph = JSON.parse(rawData);

    for (const key of Object.keys(graph.nodes)) {
      const node = graph.nodes[key];

      if (BLOCKLIST.includes(node.id)) {
        continue;
      }

      node.name = generateUniqueName();
      node.year = generateYear();
    }

    console.log(JSON.stringify(graph, null, 2));
  } catch (err) {
    console.error("Failed to process graph data:", err.message);
    process.exit(1);
  }
}

function generateUniqueName() {
  let newName;
  let safetyCounter = 0;

  do {
    const adj = ADJECTIVES[Math.floor(Math.random() * ADJECTIVES.length)];
    const noun = NOUNS[Math.floor(Math.random() * NOUNS.length)];

    newName = `${adj}${noun}`;

    safetyCounter++;
    if (safetyCounter > 1000) {
      throw new Error("Ran out of unique names! Add more adjectives or nouns.");
    }
  } while (usedFakeNames.has(newName));

  usedFakeNames.add(newName);
  return newName;
}

// Generate a random year between 2007 and 2025.
function generateYear() {
  return Math.floor(Math.random() * (2025 - 2007 + 1) + 2007);
}

run();
