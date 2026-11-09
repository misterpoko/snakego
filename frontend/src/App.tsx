import GameBoard from "./components/GameBoard";
import InputButton from "./components/InputButton";
import { useEffect, useState } from "react";
import "./App.css";

var snakeNumber = 0;

async function helloWorld() {
  const res = await fetch("/createGame");
  console.log(res);
}

async function addSnake() {
  const res = await fetch("/addSnake");
  const data = await res.json();
  console.log(data);
  snakeNumber = data.Snake ?? -1;
}

async function quitAll() {
  const res = await fetch("/quitAll");
  console.log(res);
}

async function ldir() {
  changeDir("left", snakeNumber);
}

async function rdir() {
  changeDir("right", snakeNumber);
}

async function udir() {
  changeDir("up", snakeNumber);
}

async function ddir() {
  changeDir("down", snakeNumber);
}

async function changeDir(direction: string, snake: number) {
  const res = await fetch("/changeDir", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      Snake: snake,
      Dir: direction,
    }),
  });
  console.log(res);
}

async function getInfo() {
  const res = await fetch("/gameInfo");
  return res.json();
}

function App() {
  const [gameState, setGameState] = useState<any>(null);

  useEffect(() => {
    async function poll() {
      try {
        const data = await getInfo();
        setGameState(data);
      } catch (err) {
        console.error(err);
      }
    }

    poll();
    const id = setInterval(poll, 100);
    return () => clearInterval(id);
    console.log(gameState);
  }, []);
  useEffect(() => {
    function onKey(e: KeyboardEvent) {
      if (e.key === "ArrowUp" || e.key === "w") udir();
      if (e.key === "ArrowDown" || e.key === "s") ddir();
      if (e.key === "ArrowLeft" || e.key === "a") ldir();
      if (e.key === "ArrowRight" || e.key === "d") rdir();
    }
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, []);

  return (
    <>
      <GameBoard grid={gameState?.Grid ?? null} />
      <InputButton
        text="startGame"
        className="menuButton"
        onClick={helloWorld}
      />

      <InputButton text="addPlayer" className="menuButton" onClick={addSnake} />

      <div className="dpad">
        <InputButton text="▲" className="dpadButton up" onClick={udir} />
        <InputButton text="◀" className="dpadButton left" onClick={ldir} />
        <InputButton text="▶" className="dpadButton right" onClick={rdir} />
        <InputButton text="▼" className="dpadButton down" onClick={ddir} />
      </div>

      <InputButton text="quit" className="menuButton" onClick={quitAll} />
    </>
  );
}

export default App;
