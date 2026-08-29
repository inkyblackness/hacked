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
   {
      .title = "Simple DirectMedia Layer",
      .url = "https://www.libsdl.org",
      .text = "Copyright (C) 1997-2026 Sam Lantinga <slouken@libsdl.org>\n"
              "\n"
              "This software is provided 'as-is', without any express or implied\n"
              "warranty.  In no event will the authors be held liable for any damages\n"
              "arising from the use of this software.\n"
              "\n"
              "Permission is granted to anyone to use this software for any purpose,\n"
              "including commercial applications, and to alter it and redistribute it\n"
              "freely, subject to the following restrictions:\n"
              "\n"
              "1. The origin of this software must not be misrepresented; you must not\n"
              "   claim that you wrote the original software. If you use this software\n"
              "   in a product, an acknowledgment in the product documentation would be\n"
              "   appreciated but is not required.\n"
              "2. Altered source versions must be plainly marked as such, and must not be\n"
              "   misrepresented as being the original software.\n"
              "3. This notice may not be removed or altered from any source distribution.\n",
   },
   {
      // This license also includes the C-binding, as only the output of the generator is used.
      .title = "Dear ImGui",
      .url = "https://github.com/ocornut/imgui",
      .text = "The MIT License (MIT)\n"
              "\n"
              "Copyright (c) 2014-2026 Omar Cornut\n"
              "\n"
              "Permission is hereby granted, free of charge, to any person obtaining a copy\n"
              "of this software and associated documentation files (the \"Software\"), to deal\n"
              "in the Software without restriction, including without limitation the rights\n"
              "to use, copy, modify, merge, publish, distribute, sublicense, and/or sell\n"
              "copies of the Software, and to permit persons to whom the Software is\n"
              "furnished to do so, subject to the following conditions:\n"
              "\n"
              "The above copyright notice and this permission notice shall be included in all\n"
              "copies or substantial portions of the Software.\n"
              "\n"
              "THE SOFTWARE IS PROVIDED \"AS IS\", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR\n"
              "IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,\n"
              "FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE\n"
              "AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER\n"
              "LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,\n"
              "OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE\n"
              "SOFTWARE.\n",
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
