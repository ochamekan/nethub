import type { FactoryOpts } from "imask";

export const hostnameMask: FactoryOpts = {
  mask: /^[a-zA-Z0-9._-\s]*$/,
};

export const locationMask: FactoryOpts = {
  mask: /^[a-zA-Zа-яА-ЯёЁ\s-]*$/,
};

export const ipMask: FactoryOpts = {
  mask: /^[0-9a-fA-F.:/]*$/,
};
