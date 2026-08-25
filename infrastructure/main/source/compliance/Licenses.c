#include "hacked/infrastructure/compliance/Licenses.h"

static LicenseInfo licenses[] = {
   {
      .title = "InkyBlackness - HackEd",
      .url = "https://inkyblackness.github.io",
      .text = "New BSD License\n"
              "\n"
              "© 2026, Christian Haas\n"
              "All rights reserved.\n"
              "\n"
              "Redistribution and use in source and binary forms, with or without\n"
              "modification, are permitted provided that the following conditions are met:\n"
              "\n"
              "    * Redistributions of source code must retain the above copyright\n"
              "      notice, this list of conditions and the following disclaimer.\n"
              "    * Redistributions in binary form must reproduce the above copyright\n"
              "      notice, this list of conditions and the following disclaimer in the\n"
              "      documentation and/or other materials provided with the distribution.\n"
              "    * Neither the name of \" InkyBlackness \" nor the names of its\n"
              "      contributors may be used to endorse or promote products derived from\n"
              "      this software without specific prior written permission.\n"
              "\n"
              "THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS \" AS IS \" AND\n"
              "ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED\n"
              "WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE\n"
              "DISCLAIMED. IN NO EVENT SHALL THE LISTED COPYRIGHT HOLDERS BE LIABLE FOR ANY\n"
              "DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES\n"
              "(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;\n"
              "LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND\n"
              "ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT\n"
              "(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS\n"
              "SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.\n",
   },
};

size_t licensesGetLicenseCount()
{
   return sizeof(licenses) / sizeof(LicenseInfo);
}

LicenseInfo const *licensesGetLicense(size_t const index)
{
   if (index >= licensesGetLicenseCount())
   {
      return NULL;
   }
   return &licenses[index];
}
