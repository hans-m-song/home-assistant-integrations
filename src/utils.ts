export const sleep = (timeout: number) =>
  new Promise((resolve) => setTimeout(resolve, timeout));

export const asyncInterval = (fn: () => Promise<void>, minDelay: number) => {
  console.log("beginning asyncInterval", minDelay);

  let next = true;
  const stop = () => {
    console.log("stopping asyncInterval");
    next = false;
  };

  const iterate = async () => {
    console.log("beginning iteration");

    const now = Date.now();
    await fn();
    const duration = Date.now() - now;
    const timeout = Math.max(minDelay - duration, 0);

    if (timeout > 0) {
      await sleep(timeout);
    }

    if (next) {
      await iterate();
    }
  };

  return [iterate, stop] as const;
};

export const unpackError = (error: any) => {
  const { code, name, message } = error ?? {};
  return { code, name, message };
};
