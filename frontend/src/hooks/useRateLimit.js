import { useState, useCallback, useRef } from 'react';

export function useRateLimit(cooldownSeconds = 3) {
  const [isOnCooldown, setIsOnCooldown] = useState(false);
  const [remainingTime, setRemainingTime] = useState(0);
  const timeoutRef = useRef(null);
  const intervalRef = useRef(null);

  const startCooldown = useCallback(() => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }
    if (intervalRef.current) {
      clearInterval(intervalRef.current);
    }

    setIsOnCooldown(true);
    setRemainingTime(cooldownSeconds);

    intervalRef.current = setInterval(() => {
      setRemainingTime((prev) => {
        const newTime = prev - 1;
        if (newTime <= 0) {
          clearInterval(intervalRef.current);
          return 0;
        }
        return newTime;
      });
    }, 1000);

    timeoutRef.current = setTimeout(() => {
      setIsOnCooldown(false);
      setRemainingTime(0);
      clearInterval(intervalRef.current);
    }, cooldownSeconds * 1000);
  }, [cooldownSeconds]);

  const executeWithRateLimit = useCallback(
    async (fn) => {
      if (isOnCooldown) {
        return { rateLimited: true };
      }

      startCooldown();

      try {
        const result = await fn();
        return { rateLimited: false, ...result };
      } catch (error) {
        return { rateLimited: false, error };
      }
    },
    [isOnCooldown, startCooldown]
  );

  return {
    executeWithRateLimit,
    isOnCooldown,
    remainingTime,
  };
}
