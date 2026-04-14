export type EventType =
  | "urlaub"
  | "urlaub halbtags"
  | "workation"
  | "krank"
  | "krank halbtags"
  | "home office"
  | "zeitausgleich";

const eventMap: Record<EventType, string> = {
  urlaub: "Urlaub",
  "urlaub halbtags": "Urlaub Halbtags",
  krank: "Krank",
  "krank halbtags": "Krank Halbtags",
  "home office": "Home Office",
  workation: "Workation",
  zeitausgleich: "Zeitausgleich",
};

const eventIdMap: Record<EventType, string> = {
  urlaub: "U",
  "urlaub halbtags": "UH",
  krank: "K",
  "krank halbtags": "KH",
  "home office": "HO",
  workation: "W",
  zeitausgleich: "Z",
};

export function getEventLabel(event: EventType): string {
  return eventMap[event];
}

export function getEventId(event: EventType): string {
  return eventIdMap[event];
}

export function getEventMap(): Record<EventType, string> {
  return eventMap;
}
