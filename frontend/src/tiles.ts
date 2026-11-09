export const OBJECT_TYPE_FLAG = 0b111;

export const Tile = {
  NOTHING: 0b000,
  SNAKE: 0b001,
  APPLE: 0b010,
  HALF_ME: 0b011,
  DECAY_TAIL: 0b100,
  PORTAL_ORANGE: 0b101,
  PORTAL_BLUE: 0b110,
} as const;

export type Tile = (typeof Tile)[keyof typeof Tile];

type TileStyle = { char: string; color: string };

export const TILE_STYLES: Record<Tile, TileStyle> = {
  [Tile.NOTHING]: { char: " ", color: "black" },
  [Tile.SNAKE]: { char: "■", color: "rgb(0, 255, 0)" },
  [Tile.APPLE]: { char: "●", color: "red" },
  [Tile.HALF_ME]: { char: "◐", color: "white" },
  [Tile.DECAY_TAIL]: { char: "■", color: "rgb(139, 69, 19)" },
  [Tile.PORTAL_ORANGE]: { char: "◉", color: "rgb(255, 140, 0)" },
  [Tile.PORTAL_BLUE]: { char: "◉", color: "rgb(0, 120, 255)" },
};

export function getTileType(value: number): Tile {
  return (value & OBJECT_TYPE_FLAG) as Tile;
}
