import React, { useEffect, useState } from 'react';
import '../styles/Snow.css';

function SnowEffect() {
  const [snowflakes, setSnowflakes] = useState([]);

  useEffect(() => {
    const flakes = [];
    const numberOfFlakes = 50;

    for (let i = 0; i < numberOfFlakes; i++) {
      flakes.push({
        id: i,
        left: Math.random() * 100,
        animationDuration: Math.random() * 3 + 5,
        animationDelay: Math.random() * 5,
        fontSize: Math.random() * 10 + 10,
        opacity: Math.random() * 0.6 + 0.4,
      });
    }

    setSnowflakes(flakes);
  }, []);

  return (
    <div className="snow-container">
      {snowflakes.map((flake) => (
        <div
          key={flake.id}
          className="snowflake"
          style={{
            left: `${flake.left}%`,
            animationDuration: `${flake.animationDuration}s`,
            animationDelay: `${flake.animationDelay}s`,
            fontSize: `${flake.fontSize}px`,
            opacity: flake.opacity,
          }}
        >
          ❄
        </div>
      ))}
    </div>
  );
}

export default SnowEffect;
