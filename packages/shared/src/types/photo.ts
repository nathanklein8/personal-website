export interface Photo {
  id: number;
  title: string;
  filePath: string;
  altText?: string | null;
  dateTaken?: string | null;
  location?: string | null;
  camera?: string | null;
  lens?: string | null;
  aperture?: string | null;
  shutterSpeed?: string | null;
  iso?: string | null;
  visible: boolean;
  featured: boolean;
  sortOrder: number;
  sourcePath: string;
  thumbnailPath?: string | null;
  mediumPath?: string | null;
  year: string;
  event: string;
  filename: string;
}
