# Scooter Distribution

## City Data and Bounding Boxes

- **Ottawa**
  - Population: ~1,072,000
  - Scooters: 1,072 (conservative, 1 per 1,000 inhabitants)
  - Bounding Box:
    - NorthWest: 45.37286078814053, -75.97089093158945
    - SouthEast: 45.1090266741051, -75.22792851068897

- **Montreal**
  - Population: ~1,792,000
  - Scooters: 1,792 (conservative, 1 per 1,000 inhabitants)
  - Bounding Box:
    - NorthWest: 45.62109228798646, -73.63465335011105
    - SouthEast: 45.452507945877, -73.55019903119938

## Density and Seeding Rationale

- The number of scooters is based on a low, conservative density: 1 scooter per 1,000 inhabitants.
- Scooters are seeded randomly within the bounding box for each city.

## Geographical Limitations

- Bounding boxes are a simplification: real deployment areas are irregular and may include unpopulated or suboptimal regions.
- Making the box too small excludes optimal areas; making it too large includes depopulated or non-usable zones.
- For this exercise, a rectangular bounding box is used for simplicity and reproducibility.
