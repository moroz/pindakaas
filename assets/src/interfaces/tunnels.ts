import type { ISOTimestamp, UUID } from "./common";

export interface TunnelListDTO {
  id: UUID;
  subdomain: string;
  username: string;
  insertedAt: ISOTimestamp;
  updatedAt: ISOTimestamp;
  userId: UUID;
  active: boolean;
}

export interface TunnelListResponse {
  data: ReadonlyArray<TunnelListDTO>;
}
