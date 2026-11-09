import "./GameBoard.css";
import { TILE_STYLES, getTileType } from "../tiles";

export type Grid = number[][];

function GameBoard({ grid }: { grid: Grid | null }) {
  if (!grid) return <div>Loading...</div>;

  const cells = [];
  for (let y = 0; y < grid.length; y++) {
    for (let x = 0; x < grid[y].length; x++) {
      const style = TILE_STYLES[getTileType(grid[y][x])];
      cells.push(
        <div key={`${x}-${y}`} className="cell" style={{ color: style.color }}>
          {style.char}
        </div>,
      );
    }
  }

  return <div className="board">{cells}</div>;
}

export default GameBoard;
