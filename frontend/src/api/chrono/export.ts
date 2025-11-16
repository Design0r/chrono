import { CHRONO_URL } from "./chrono";

export class ApiExport {
  async download(year: number): Promise<void> {
    const a = document.createElement("a");
    a.href = `${CHRONO_URL}/export/${year}`;
    a.download = `chrono_krankheitstage_export_${year}.csv`;
    a.click();
  }
}
