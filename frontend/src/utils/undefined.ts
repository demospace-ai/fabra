export function toNull<T>(arg: T | undefined): T | null {
  if (arg === undefined) {
    return null;
  }

  return arg;
}

export function toEmptyList<T>(arg: T | undefined): T | [] {
  if (arg === undefined || arg === null) {
    return [];
  }

  return arg;
}

export function toUndefined<T>(arg: T | null): T | undefined {
  if (arg === null) {
    return undefined;
  }

  return arg;
}
