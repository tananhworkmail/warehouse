// Pair capacities supplied in "PALET KTP.xlsx" (sheet 3.2026).
// Grouping rack codes by capacity keeps the source data compact and auditable.
const rackCodesByCapacity = {
  1008: ["A43", "B57"],
  1260: ["A02", "A06", "A39", "A49", "A53", "B49"],
  1512: [
    "A01", "A03", "A04", "A05", "A07", "A08", "A09", "A10", "A11",
    "A48", "A50", "A51", "A52", "A54", "A55", "A56", "A57", "A58",
    "B45", "B46", "B47", "B48", "B50", "B51", "B52", "B53", "B54",
    "B55", "B56",
  ],
  1764: [
    "A29", "A30", "A31", "A32", "A33", "A34", "A35", "A36", "A37",
    "A38", "A40", "A41", "A42", "A44", "A45", "A46", "A47",
  ],
  2016: ["B62", "B67"],
  2520: [
    "A12", "B58", "B59", "B60", "B61", "B63", "B64", "B65", "B66",
    "B68", "B69", "B70", "B71",
  ],
  2772: ["B33"],
  3024: ["B09"],
  3276: ["A19", "A23", "A27"],
  3528: ["B05", "B06", "B07", "B08"],
  3780: [
    "B14", "B34", "B35", "B36", "B37", "B38", "B39", "B40", "B41",
    "B42", "B43", "B44",
  ],
  4284: [
    "A13", "A14", "A15", "A16", "A17", "A18", "A20", "A21", "A22",
    "A24", "A25", "A26", "A28",
  ],
  4536: ["B10", "B11", "B12", "B13", "B23", "B28"],
  4788: ["B01", "B02", "B03", "B04"],
  5292: [
    "B15", "B16", "B17", "B18", "B19", "B20", "B21", "B22", "B24",
    "B25", "B26", "B27", "B29", "B30", "B31", "B32",
  ],
};

export const ktpRackCapacityByCode = Object.freeze(
  Object.fromEntries(
    Object.entries(rackCodesByCapacity).flatMap(([capacity, rackCodes]) =>
      rackCodes.map((rackCode) => [rackCode, Number(capacity)]),
    ),
  ),
);

export const ktpWarehouseCapacity = Object.values(
  ktpRackCapacityByCode,
).reduce((total, capacity) => total + capacity, 0);
